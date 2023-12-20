import { API } from './store/api';

export default class Store {

    constructor( DI ) {
        DI.bind( 'store', this );
        this._setupAPI( DI.resolve( 'api:url' ), DI.resolve( 'api:headers' ) );
    }

    _setupAPI( baseURL, headers ) {
        this.api = new API( {
            baseURL,
            headers
        } );
    }

    articles( offset, limit, rest = {} ) {
        return this.api.request( {
            url: '/articles',
            params: {
                offset,
                limit,
                ...rest
            }
        } );
    }

    articleBySlug( slug ) {
        return this.api.request( {
            url: `/articles/${slug}`
        } );
    }
}
