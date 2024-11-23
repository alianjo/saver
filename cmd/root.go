package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/alianjo/saver/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getLineNumber retrieves the line number of the caller
func getLineNumber() int {
	_, _, line, _ := runtime.Caller(1) // The third return value is the line number
	return line
}

var opt struct {
	namespace string
	output    string
}

// NewBackupCMD creates a new Cobra command for backup
func NewBackupCMD() *cobra.Command {
	rootcmd := &cobra.Command{
		Use:     "save",
		Short:   "sv",
		Long:    "save workloads",
		Example: "kubectl save deployment --namespace controller -o controller-deployments.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// Correctly report the file and line number
				return fmt.Errorf("Pass one of the options [deployment, statefulset, daemonset] is required at %s:%d", "cmd/backup.go", getLineNumber())
			}
			// Use the package-level opt struct for namespace and workload
			return backup(args[0], opt.namespace, opt.output)
		},
		Version: "0.1",
	}

	// Correctly bind the flags to the opt struct
	rootcmd.Flags().StringVarP(&opt.namespace, "namespace", "n", "default", " namespace")
	rootcmd.Flags().StringVarP(&opt.output, "output", "o", "", "output path")

	return rootcmd
}

// backup simulates the backup process for a Kubernetes secret

func PrintWorkloadInfo(clientset *kubernetes.Clientset, namespace string, workload string) {
	ctx := context.Background()

	switch workload {
	case "deployment":
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching deployments: %v\n", err)
			return
		}
		fmt.Printf("Namespace: %s | Deployments:\n", namespace)
		for _, d := range deployments.Items {
			fmt.Printf("  - %s (Replicas: %d)\n", d.Name, *d.Spec.Replicas)
		}
	case "daemonset":
		daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching daemonsets: %v\n", err)
			return
		}
		fmt.Printf("Namespace: %s | DaemonSets:\n", namespace)
		for _, ds := range daemonSets.Items {
			fmt.Printf("  - %s (Desired Pods: %d)\n", ds.Name, ds.Status.DesiredNumberScheduled)
		}
	case "statefulset":
		statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching statefulsets: %v\n", err)
			return
		}
		fmt.Printf("Namespace: %s | StatefulSets:\n", namespace)
		for _, ss := range statefulSets.Items {
			fmt.Printf("  - %s (Replicas: %d)\n", ss.Name, *ss.Spec.Replicas)
		}
	default:
		fmt.Printf("Invalid workload type: %s\n", workload)
	}
}

func PrintWorkloadYaml(clientset *kubernetes.Clientset, namespace string, workload string) {
	ctx := context.Background()

	switch workload {
	case "deployment":
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching deployments: %v\n", err)
			return
		}
		for _, d := range deployments.Items {
			printYaml(d)
		}
	case "daemonset":
		daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching daemonsets: %v\n", err)
			return
		}
		for _, ds := range daemonSets.Items {
			printYaml(ds)
		}
	case "statefulset":
		statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching statefulsets: %v\n", err)
			return
		}
		for _, ss := range statefulSets.Items {
			printYaml(ss)
		}
	default:
		fmt.Printf("Invalid workload type: %s\n", workload)
	}
}

// PrintWorkloadYamlToFile fetches the YAML of the specified workload resources and writes them to a file.
func PrintWorkloadYamlToFile(clientset *kubernetes.Clientset, namespace string, workload string, filename string) error {
	ctx := context.Background()

	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	switch workload {
	case "deployment":
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching deployments: %v", err)
		}
		for _, d := range deployments.Items {
			if err := writeYamlToFile(file, d); err != nil {
				return err
			}
		}
	case "daemonset":
		daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching daemonsets: %v", err)
		}
		for _, ds := range daemonSets.Items {
			if err := writeYamlToFile(file, ds); err != nil {
				return err
			}
		}
	case "statefulset":
		statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching statefulsets: %v", err)
		}
		for _, ss := range statefulSets.Items {
			if err := writeYamlToFile(file, ss); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid workload type: %s", workload)
	}

	return nil
}

func writeYamlToFile(file *os.File, resource interface{}) error {
	yamlData, err := yaml.Marshal(resource)
	if err != nil {
		return fmt.Errorf("error marshalling resource to YAML: %v", err)
	}
	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}
	_, err = file.WriteString("\n---\n") // Separate resources in the file with ---
	return err
}

func printYaml(resource interface{}) {
	yamlData, err := yaml.Marshal(resource)
	if err != nil {
		fmt.Printf("Error marshalling resource to YAML: %v\n", err)
		return
	}
	fmt.Println("---")
	fmt.Println(string(yamlData))
}

func backup(workload, namespace, output string) error {
	fmt.Printf("Backing up %s  in namespace %s\n", workload, namespace)
	client := client.K8Client()

	PrintWorkloadInfo(client, namespace, workload)
	PrintWorkloadYaml(client, namespace, workload)
	if output != "" {
		PrintWorkloadYamlToFile(client, namespace, workload, output)
	}

	return nil
}
