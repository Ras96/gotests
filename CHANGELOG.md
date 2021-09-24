# CHANGELOG

## Changes

[🎉 Use sprig.Funcs() · Ras96/gotests@3eab362](https://github.com/Ras96/gotests/commit/3eab362767997cfee218a13a512a5106a9b6689f)

- Use `html/template` instead of `text/template`
- Use `github.com/Masterminds/sprig`
- Update go.mod & go.sum

```diff
diff --git a/internal/render/render.go b/internal/render/render.go
index 4317420..a6f5aa1 100755
--- a/internal/render/render.go
+++ b/internal/render/render.go
@@ -5,11 +5,12 @@ import (
        "io"
        "io/ioutil"
        "path"
-       "text/template"
+       "html/template"

        "github.com/cweill/gotests/internal/models"
        "github.com/cweill/gotests/internal/render/bindata"
        "github.com/cweill/gotests/templates"
+       "github.com/Masterminds/sprig"
 )

 type Render struct {
@@ -24,7 +25,7 @@ func New() *Render {
                        "Param":    parameterName,
                        "Want":     wantName,
                        "Got":      gotName,
-               }),
+               }).Funcs(sprig.FuncMap()),
        }

        // default templates first
```
