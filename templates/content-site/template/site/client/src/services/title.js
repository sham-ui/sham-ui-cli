export default class Title {
    constructor( DI ) {
        DI.bind( 'title', this );
        this.document = DI.resolve( 'document' );
    }

    /**
     * Set new document title
     * @param {string} newTitle
     */
    change( newTitle ) {
        this.document.title = `{{ logoText }} | ${newTitle}`;
    }
}
