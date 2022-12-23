;; Learn Fennel: fennel-lang.org/tutorial
;; Style Guide: fennel-lang.org/style

(local module {:layers {}})

(fn module.search [query progress]
  (progress (.. "Searching for " query))
  (let [medias {}]
    medias))

(fn layer [media progress]
  (progress (.. "Layer " media.display))
  (let [medias {}]
    medias))

(tset module.layers "First Layer" layer)

(fn module.prepare [media progress]
  (progress (.. "Preparing " media.display))
  media)

(fn module.stream [media progress]
  (progress (.. "Streaming " media.display)))

(fn module.download [media progress]
  (progress (.. "Downloading " media.display)))

module

