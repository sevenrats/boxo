## Usage

Here's how you create, start, interact with, and stop the provider system:

```golang
import (
	"context"
	"time"

	"github.com/sevenrats/boxo/provider"
	"github.com/sevenrats/boxo/provider/queue"
	"github.com/sevenrats/boxo/provider/simple"
)

rsys := (your routing system here)
dstore := (your datastore here)
cid := (your cid to provide here)

q := queue.NewQueue(context.Background(), "example", dstore)

reprov := simple.NewReprovider(context.Background(), time.Hour * 12, rsys, simple.NewBlockstoreProvider(dstore))
prov := simple.NewProvider(context.Background(), q, rsys)
sys := provider.NewSystem(prov, reprov)

sys.Run()

sys.Provide(cid)

sys.Close()
```
