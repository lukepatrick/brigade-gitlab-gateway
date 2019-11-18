package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brigadecore/brigade/pkg/brigade"

	whgitlab "github.com/brigadecore/brigade-gitlab-gateway/pkg/webhook"

	"k8s.io/api/core/v1"

	"github.com/brigadecore/brigade/pkg/storage/kube"

	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
)

var (
	kubeconfig  string
	master      string
	namespace   string
	gatewayPort string
)

const (
	path = "/events/gitlab"
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.StringVar(&namespace, "namespace", defaultNamespace(), "kubernetes namespace")
	flag.StringVar(&gatewayPort, "gateway-port", defaultGatewayPort(), "TCP port to use for brigade-gitlab-gateway")
}

func main() {
	flag.Parse()

	hook := gitlab.New(&gitlab.Config{Secret: ""})
	hook.RegisterEvents(HandleMultiple,
		gitlab.PushEvents,
		gitlab.TagEvents,
		gitlab.IssuesEvents,
		gitlab.ConfidentialIssuesEvents,
		gitlab.CommentEvents,
		gitlab.MergeRequestEvents,
		gitlab.WikiPageEvents,
		gitlab.PipelineEvents,
		gitlab.BuildEvents) // Add as many as needed

	err := webhooks.Run(hook, ":"+gatewayPort, path)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func defaultNamespace() string {
	if ns, ok := os.LookupEnv("BRIGADE_NAMESPACE"); ok {
		return ns
	}
	return v1.NamespaceDefault
}

func defaultGatewayPort() string {
	if port, ok := os.LookupEnv("BRIGADE_GITLAB_GATEWAY_PORT"); ok {
		return port
	}
	return "7746"
}

// HandleMultiple handles multiple GitLab events
func HandleMultiple(payload interface{}, header webhooks.Header) {
	log.Println("Handling Payload..")

	clientset, err := kube.GetClient(master, kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	store := kube.New(clientset, namespace)
	store.GetProjects()

	glhandler := whgitlab.NewGitlabHandler(store)

	var repo, secret string
	var rev brigade.Revision
	secret = strings.Join(header["X-Gitlab-Token"], "")

	switch payload.(type) {
	case gitlab.PushEventPayload:
		log.Println("case gitlab.PushEventPayload")
		release := payload.(gitlab.PushEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Commit = release.CheckoutSHA
		rev.Ref = release.Ref

		glhandler.HandleEvent(repo, "push", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.TagEventPayload:
		log.Println("case gitlab.TagEventPayload")
		release := payload.(gitlab.TagEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Commit = release.CheckoutSHA
		rev.Ref = release.Ref

		glhandler.HandleEvent(repo, "tag", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.IssueEventPayload:
		log.Println("case gitlab.IssueEventPayload")
		release := payload.(gitlab.IssueEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Ref = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "issue", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.ConfidentialIssueEventPayload:
		log.Println("case gitlab.ConfidentialIssueEventPayload")
		release := payload.(gitlab.ConfidentialIssueEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Ref = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "issue", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.CommentEventPayload:
		log.Println("case gitlab.CommentEventPayload")
		release := payload.(gitlab.CommentEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Commit = release.Commit.ID

		glhandler.HandleEvent(repo, "comment", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.MergeRequestEventPayload:
		log.Println("case gitlab.MergeRequestEventPayload")
		release := payload.(gitlab.MergeRequestEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Ref = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "mergerequest", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.WikiPageEventPayload:
		log.Println("case gitlab.WikiPageEventPayload")
		release := payload.(gitlab.WikiPageEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Ref = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "wikipage", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.PipelineEventPayload:
		log.Println("case gitlab.PipelineEventPayload")
		release := payload.(gitlab.PipelineEventPayload)

		repo = release.Project.PathWithNamespace
		rev.Commit = release.Commit.ID

		glhandler.HandleEvent(repo, "pipeline", rev, []byte(fmt.Sprintf("%v", release)), secret)

	case gitlab.BuildEventPayload:
		log.Println("case gitlab.BuildEventPayload")
		release := payload.(gitlab.BuildEventPayload)

		repo = release.ProjectName
		rev.Commit = release.SHA
		rev.Ref = release.Ref

		glhandler.HandleEvent(repo, "build", rev, []byte(fmt.Sprintf("%v", release)), secret)

	default:
		log.Printf("Unsupported event")
		return
	}

}
