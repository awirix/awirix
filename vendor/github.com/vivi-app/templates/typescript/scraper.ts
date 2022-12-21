/**
 * Learn More: typescripttolua.github.io
 * Caveats:    typescripttolua.github.io/docs/caveats
 * Compile:    `npm install && npm run build`
 */

/**
 * Intermediate value that is passed between the states
 */
interface Media {
  /**
   * Used for the string representation inside vivi's interface.
   */
  display: string

  /**
   * Optional field that, if present, will be used
   * as a short description of the media.
   */
  about?: string

  // Any other fields will be preserved when passing media
  // between states. Feel free to add them, e.g. `url` field
}

/**
 * Function that will pass message to the vivi's interface
 * during loading.
 */
type Progress = (message: string) => void

/**
 * Searches for the media
 * This function may be omitted if this extension does not provide searching functionality.
 * Might be the case if it is dedicated to the single show/movie/book/...
 */
export function search(query: string, progress: Progress): Media[] {
  progress(`Searching for ${query}`)

  return []
}

function layer(media: Media, progress: Progress): Media[] {
  progress(`Layer ${media.display}`)

  return []
}

/**
 * Each layer returns a list of sub-media for the given one.
 * For example, you can search for a show, then selected show will be passed to the first layer that's responsible for returning show's seasons.
 * After that, the selected season will be passed to the second layer that would return season's episodes.
 * Layers may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just search and watch, no seasons, no episodes).
 * If `search` function is omitted, first layer will receive a `null` instead of the media.
 */
export const layers = {
  "First Layer": layer
} as {
  [name: string]: (media: Media, progress: Progress) => Media[]
}

/**
 * Prepares a media for streaming/downloading.
 * When you need to do some heavy operations before passing it further.
 * E.g. decode a video stream url, calculate a hash of the file, etc.
 */
export function prepare(media: Media, progress: Progress): Media {
  progress(`Preparing ${media.display}`)

  return media
}

/**
 * Streams a media
 * May be omitted if the extension does not provide streaming functionality.
 * However, at least one of the `stream` or `download` functions must be present.
 */
export function stream(media: Media, progress: Progress) {
  progress(`Streaming ${media.display}`)

  // @ts-ignore
  error('Not implemented')
}

/**
 * Downloads a media.
 * May be omitted if this extension does not provide downloading functionality.
 * However, at least one of the `stream` or `download` functions must be present.
 */
export function download(media: Media, progress: Progress) {
  progress(`Downloading ${media.display}`)

  // @ts-ignore
  error('Not implemented')
}
