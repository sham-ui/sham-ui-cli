export default class Cookies {
    constructor( DI ) {
        DI.bind( 'cookies', this );
        this.document = DI.resolve( 'document' );
    }

    /**
     * @param {string} name
     * @return {string|undefined}
     */
    get( name ) {
        const escaped = name.replace( /([.$?*|{}()[\]\\/+^])/g, '\\$1' );
        const matches = this.document.cookie.match(
            new RegExp(
                '(?:^|; )' + escaped + '=([^;]*)'
            )
        );
        return matches ?
            matches[ 1 ] :
            undefined;
    }

    /**
     * @param {string} name
     * @param {string} value
     * @param {Object} options
     */
    set( name, value, options = {} ) {
        options = {
            path: '/',
            ...options
        };

        if ( options.expires instanceof Date ) {
            options.expires = options.expires.toUTCString();
        }

        let updatedCookie = `${name}=${value}`;

        for ( let optionKey in options ) {
            updatedCookie += '; ' + optionKey;
            const optionValue = options[ optionKey ];
            if ( optionValue !== true ) {
                updatedCookie += '=' + optionValue;
            }
        }

        this.document.cookie = updatedCookie;
    }

    /**
     * @param {string} name
     */
    delete( name ) {
        this.set( name, '', { 'max-age': -1 } );
    }
}
