interface Media {
  display: string
  about?: string
  url?: string
}

type Progress = (message: string) => void

export function {{ .Fn.Search }}(query: string, progress: Progress): Media[] {
  progress(`Searching for ${query}`)

  return []
}

function layer(media: Media, progress: Progress): Media[] {
  progress(`Layer ${media.display}`)

  return []
}

export const {{ .Fields.Layers }} = {
  "First Layer": layer
} as {
  [name: string]: (media: Media, progress: Progress) => Media[]
}

export function {{ .Fn.Prepare }}(media: Media, progress: Progress): Media {
  progress(`Preparing ${media.display}`)

  return media
}

export function {{ .Fn.Stream }}(media: Media, progress: Progress) {
  progress(`Streaming ${media.display}`)

  error('Not implemented')
}

export function {{ .Fn.Download }}(media: Media, progress: Progress) {
  progress(`Downloading ${media.display}`)

  error('Not implemented')
}
