;; Learn Fennel: fennel-lang.org/tutorial
;; Style Guide: fennel-lang.org/style

(local module {})

(local search {:handler (fn [query progress]
                          [])})

(local layers [{:handler (fn [media progress]
                           [])}])

(tset module :search search)
(tset module :layers layers)

(fn module.prepare [media progress]
  media)

(fn module.stream [media progress]
  (error "Not implemented"))

(fn module.download [media progress]
  (error "Not implemented"))

module

