export default class Title {
    constructor( DI ) {
        DI.bind( 'title', this );
    }

    /**
     * Set new document title
     * @param {string} newTitle
     */
    change( newTitle ) {
        document.title = `{{ logoText }} | ${newTitle}`;
    }
}
