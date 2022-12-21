(local module {})

(fn module.test []
  (assert (= (+ 2 2) 4) "Math is broken"))

module
