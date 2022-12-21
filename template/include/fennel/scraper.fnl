(local module {:layers {}})

(fn module.{{ .Fn.Search }} [query progress]
  (progress (.. "Searching for " query))
  (let [medias {}]
    medias))

(fn layer [media progress]
  (progress (.. "Layer " media.display))
  (let [medias {}]
    medias))

(tset module.{{ .Fields.Layers }} "First Layer" layer)

(fn module.{{ .Fn.Prepare }} [media progress]
  (progress (.. "Preparing " media.{{ .Fields.Display }}))
  media)

(fn module.{{ .Fn.Stream }} [media progress]
  (progress (.. "Streaming " media.{{ .Fields.Display }})))

(fn module.{{ .Fn.Download }} [media progress]
  (progress (.. "Downloading " media.{{ .Fields.Display }})))

module

