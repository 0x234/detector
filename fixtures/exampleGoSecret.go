package exampleSecret

import "fmt"

// Deliberately insecure Go source code to verify the detection
// rule is working.

func connectToAws() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("20664672", "hunter2", "89y32f9889dsfsd98h89f"),
	})
}

func main() {
	fmt.Printf("Let's connect to AWS")
}
