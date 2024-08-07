package letme

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	utils "github.com/lockedinspace/letme/pkg"
	"github.com/spf13/cobra"
)

var obtainCmd = &cobra.Command{
	Use:     "obtain",
	Aliases: []string{"ob"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ConfigFileHealth()
	},
	Short: "Obtain account credentials.",
	Long: `Obtain AWS STS assumed credentials once the user authenticates itself.
Credentials will last 3600 seconds by default and can be used with the argument '--profile $ACCOUNT_NAME'
within the AWS cli binary.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// get flags
		inlineTokenMfa, _ := cmd.Flags().GetString("inline-mfa")
		renew, _ := cmd.Flags().GetBool("renew")
		credentialProcess, _ := cmd.Flags().GetBool("credential-process")
		localCredentialProcessFlagV1, _ := cmd.Flags().GetBool("v1")

		// get the current context
		currentContext := utils.GetCurrentContext()
		letmeContext := utils.GetContextData(currentContext)
		if letmeContext.AwsSessionDuration == 0 {
			letmeContext.AwsSessionDuration = 3600
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(letmeContext.AwsSourceProfile), config.WithRegion(letmeContext.AwsSourceProfileRegion))
		utils.CheckAndReturnError(err)
		account := utils.GetAccount(letmeContext.AwsDynamoDbTable, cfg, args[0])

		switch {
		case len(account.Name) == 0:
			fmt.Println("letme: the specified account does not exist in your DynamoDB.")
			fmt.Println("letme: run 'letme list' to list available accounts.")
			os.Exit(1)
		case len(account.Region) == 0:
			fmt.Println("letme: default region not set. Setting 'us-east-1' by default.")
			account.Region[0] = "us-east-1"
		case len(account.Role) == 0:
			fmt.Println("letme: the specified account does not have any role configured. Nothing to assume.")
			os.Exit(1)
		}

		if credentialProcess {
			utils.AwsConfigFileCredentialsProcessV1(args[0], account.Region[0])
		}

		// overwrite the session name variable if the user provides it
		if len(letmeContext.AwsSessionName) == 0 && !localCredentialProcessFlagV1 {
			fmt.Println("Using default session name: '" + args[0] + "-letme-session' with context: '" + currentContext + "'")
			letmeContext.AwsSessionName = args[0] + "-letme-session"
		} else if !localCredentialProcessFlagV1 {
			fmt.Println("Assuming role with the following session name: '" + letmeContext.AwsSessionName + "' and context: '" + currentContext + "'")
		}

		// grab the mfa arn from the config, create a new aws session and try to get credentials
		var authMethod string
		if len(letmeContext.AwsMfaArn) > 0 && !localCredentialProcessFlagV1 {
			authMethod = "mfa"
		} else if len(letmeContext.AwsMfaArn) > 0 && localCredentialProcessFlagV1 {
			authMethod = "mfa-credential-process-v1"
		} else if localCredentialProcessFlagV1 {
			authMethod = "credential-process-v1"
		} else {
			authMethod = "assume-role"
		}

		var profileCredential utils.ProfileCredential
		var profileConfig utils.ProfileConfig
		switch {
		case len(account.Role) > 1:
			profileCredential, profileConfig = utils.AssumeRoleChained(letmeContext, cfg, inlineTokenMfa, account, renew, localCredentialProcessFlagV1, authMethod)
		default:
			profileCredential, profileConfig = utils.AssumeRole(letmeContext, cfg, inlineTokenMfa, account, renew, localCredentialProcessFlagV1, authMethod)
		}

		utils.LoadAwsCredentials(account.Name, profileCredential)
		utils.LoadAwsConfig(account.Name, profileConfig)
		fmt.Println("letme: use the argument '--profile " + account.Name + "' to interact with the account.")
	},
}

func init() {
	var credentialProcess bool
	var v1 bool
	var renew bool
	RootCmd.AddCommand(obtainCmd)
	obtainCmd.Flags().String("inline-mfa", "", "pass the mfa token without user prompt")
	obtainCmd.Flags().BoolVarP(&renew, "renew", "", false, "force new credentials to be assumed")
	obtainCmd.Flags().BoolVarP(&credentialProcess, "credential-process", "", false, "obtain credentials using the credential_process entry in your aws config file.")
	obtainCmd.Flags().BoolVarP(&v1, "v1", "", false, "output credentials following the credential_process version 1 standard.")

}
