export default class SEO {
    constructor( DI ) {
        DI.bind( 'seo', this );
        this.document = DI.resolve( 'document' );
    }

    setContent( content ) {
        this.document.content = content;
    }
}
