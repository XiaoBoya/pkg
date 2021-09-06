package route

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/katanomi/pkg/plugin/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestRepoFileGet(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestRepoFileCreator{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	httpRequest, _ := http.NewRequest("GET", "/plugins/v1alpha1/test-b/projects/1/coderepositories/1/contents/", nil)
	httpRequest.Header.Set("Accept", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	obj := metav1alpha1.GitRepoFile{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &obj)
	g.Expect(err).To(BeNil())
	g.Expect(obj.Name).To(Equal("aaa"))
}

func TestRepoFileCreate(t *testing.T) {
	g := NewGomegaWithT(t)

	ws, err := NewService(&TestRepoFileCreator{}, client.MetaFilter)
	g.Expect(err).To(BeNil())

	container := restful.NewContainer()

	container.Add(ws)

	payload := metav1alpha1.CreateRepoFileParams{
		Branch: "master",
		Message: "test",
		Content: []byte("aaaa"),
	}
	content, _ := json.Marshal(payload)
	path := "/a/aa/a/a.txt"
	filePath :=	strings.TrimLeft(path, "/")
	filePath = strings.Replace(url.PathEscape(filePath), ".", "%2E", -1)
	httpRequest, _ := http.NewRequest("POST", "/plugins/v1alpha1/test-b/projects/gitlab/coderepositories/gitlab/contents", bytes.NewBuffer(content))
	httpRequest.Header.Set("Content-type", "application/json")

	metaData := client.Meta{BaseURL: "http://api.test", Version: "v1"}
	data, _ := json.Marshal(metaData)
	meta := base64.StdEncoding.EncodeToString(data)
	httpRequest.Header.Set(client.PluginMetaHeader, meta)

	httpWriter := httptest.NewRecorder()

	container.Dispatch(httpWriter, httpRequest)
	g.Expect(httpWriter.Code).To(Equal(http.StatusOK))

	obj := metav1alpha1.GitCommit{}
	err = json.Unmarshal(httpWriter.Body.Bytes(), &obj)
	g.Expect(err).To(BeNil())
	g.Expect(obj.Spec.Message).To(Equal(&payload.Message))
}

type TestRepoFileCreator struct {
}

func (t *TestRepoFileCreator) Path() string {
	return "test-b"
}

func (t *TestRepoFileCreator) Setup(_ context.Context, _ *zap.SugaredLogger) error {
	return nil
}

func (t *TestRepoFileCreator) CreateGitRepoFile(ctx context.Context, payload metav1alpha1.CreateRepoFilePayload) (metav1alpha1.GitCommit, error) {
	return metav1alpha1.GitCommit{
		Spec: metav1alpha1.GitCommitSpec{
			Message: &payload.Message,
		},
	}, nil
}

func (t *TestRepoFileCreator) GetGitRepoFile(ctx context.Context, option metav1alpha1.GitRepoFileOption) (metav1alpha1.GitRepoFile, error) {
	return metav1alpha1.GitRepoFile{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aaa",
		},
	}, nil
}

