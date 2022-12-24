/**
 * Learn More: https://typescripttolua.github.io
 * Caveats:    https://typescripttolua.github.io/docs/caveats
 */

/**
 * Noun is used to describe the media.
 */
interface Noun {
    /**
     * Singular form of the noun
     * E.g. "Season", "Episode", "Book", etc.
     */
    singular: string;

    /**
     * Plural form of the noun
     * E.g. "Seasons", "Episodes", "Books", etc.
     */
    plural: string;
}

/**
 * Intermediate values that is passed between states.
 */
interface Media {
    /**
     * Used for the string representation inside vivi's interface.
     */
    title: string;

    /**
     * Optional field that, if present, will be used
     * as a short description of the media.
     */
    description?: string;

    // Any other fields will be preserved
    // when passing media between states.
}

interface Search {
    /**
     * Title to set in the search field.
     */
    title: string;

    /**
     * Function that will be called to search for the media.
     */
    handler: (query: string, progress: Progress) => Media[];

    /**
     * Optional field that, if present, will be used
     * as a name of the noun that represents the sub-media.
     */
    noun?: Noun;
}

interface Layer {
    /**
     * Name of the layer
     * Will be used as a title of the layer in the vivi's interface.
     * E.g. "Seasons", "Episodes", "Books", etc.
     */
    title: string;

    /**
     * Function that will be called to get the list of sub-media for the given one.
     */
    handler: (media: Media, progress: Progress) => Media[];

    /**
     * Optional field that, if present, will be used
     * as a name of the noun that represents the sub-media.
     */
    noun?: Noun;
}

/**
 * Function that will pass message to the vivi's interface
 * during loading.
 */
type Progress = (message: string) => void;

export const search = {
    title: "Search",
    handler: (query: string, progress: Progress) => [],
} as Search;

/**
 * Each layer returns a list of sub-media for the given one.
 * For example, you can search for a show, then selected show will be passed to the first layer that's responsible for returning show's seasons.
 * After that, the selected season will be passed to the second layer that would return season's episodes.
 * Layers may be omitted (nil or 0 length) if this extension does not provide such functionality (e.g. just search and watch, no seasons, no episodes).
 * If `search` function is omitted, first layer will receive a `null` instead of the media.
 */
export const layers = [
    {
        title: "First Layer",
        handler: (media: Media, progress: Progress) => [],
    },
] as Layer[];

/**
 * Prepares a media for streaming/downloading.
 * When you need to do some heavy operations before passing it further.
 * E.g. decode a video stream url, calculate a hash of the file, etc.
 */
export function prepare(media: Media, progress: Progress): Media {
    return media;
}

/**
 * Streams a media
 * May be omitted if the extension does not provide streaming functionality.
 * However, at least one of the `stream` or `download` functions must be present.
 */
export function stream(media: Media, progress: Progress) {
    // @ts-ignore
    error("Not implemented");
}

/**
 * Downloads a media.
 * May be omitted if this extension does not provide downloading functionality.
 * However, at least one of the `stream` or `download` functions must be present.
 */
export function download(media: Media, progress: Progress) {
    // @ts-ignore
    error("Not implemented");
}
