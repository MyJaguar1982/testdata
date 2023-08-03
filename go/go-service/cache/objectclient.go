// object client for cache service implementation
// The above code defines a RESTClient type that implements methods for performing GET, LIST, and
// DELETE operations on Kubernetes resources.
// @property  - - The `package objectclient` statement indicates that this code belongs to the
// `objectclient` package.

package objectclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/restclient"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
)

type RESTClient struct {
	*restclient.RESTClient
}

func (c *RESTClient) Get(ctx context.Context, key string, obj runtime.Object) error {
	return c.RESTClient.Get().
		Namespace(api.NamespaceValue(ctx)).
		Resource(util.ObjectKey(key)).
		VersionedParams(&api.GetOptions{}, api.ParameterCodec).
		Do(ctx).
		Get()
}

func (c *RESTClient) List(ctx context.Context, key string, obj runtime.Object) error {
	return c.RESTClient.Get().
		Namespace(api.NamespaceValue(ctx)).
		Resource(util.ObjectKey(key)).
		VersionedParams(&api.ListOptions{}, api.ParameterCodec).
		Do(ctx).
		Get()
}

func (c *RESTClient) Delete(ctx context.Context, key string, obj runtime.Object) error {
	return c.RESTClient.Delete().
		Namespace(api.NamespaceValue(ctx)).
		Resource(util.ObjectKey(key)).
		Do(ctx).
		Error()
}
