package ini

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestConcurrencyStress(t *testing.T) {
	ini, _ := New()
	numG := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	wg.Add(numG)
	d := t.TempDir()

	for i := 0; i < numG; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				k := fmt.Sprintf("k_%d_%d", id, j)
				ini.Set("sec", k, "v")
				_ = ini.GetString("sec", k)
				_ = ini.GetString("sec", k)
				// occasionally save
				if j%50 == 0 {
					fn := filepath.Join(d, fmt.Sprintf("out_%d_%d.ini", id, j))
					_ = ini.Save(fn, false)
					// remove it
					_ = os.Remove(fn)
				}
			}
		}(i)
	}

	wg.Wait()
}
