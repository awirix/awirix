(local module {})

(fn module.{{ .Fn.Test }} []
  (assert (= (+ 2 2) 4) "Math is broken"))

module