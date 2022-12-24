;; Learn Fennel: fennel-lang.org/tutorial
;; Style Guide: fennel-lang.org/style

(local module {})

(local search {:title :Search
               :handler (fn [query progress]
                          [])})

(local layers [{:title "First Layer"
                :handler (fn [media progress]
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

