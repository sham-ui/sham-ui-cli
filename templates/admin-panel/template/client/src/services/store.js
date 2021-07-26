import { inject } from 'sham-ui-macro/inject.macro';
import { API } from './store/api';

const VALID_SESSION_URL = '/validsession';

export default class Store {
    @inject session;

    constructor( DI ) {
        this.DI = DI;
        DI.bind( 'store', this );
        this._setupAPI();
    }

    _setupAPI() {
        const baseURL = PRODUCTION ?
            `${document.location.protocol}//${document.location.host}/api/` :
            'http://localhost:3001/api/';
        this.api = new API( this.DI, {
            baseURL,
            onUnauthorized: ::this._onAPIUnauthorized
        } );
    }

    _onAPIUnauthorized( { url } ) {
        if ( url !== VALID_SESSION_URL ) {

            // Logout if server response with 401 for any request,
            // exclude VALID_SESSION_URL
            requestAnimationFrame(
                () => this.session.logout()
            );
        }
    }

    validSession() {
        return this.api.request( { url: VALID_SESSION_URL } );
    }

    csrftoken() {
        return this.api.request( { url: '/csrftoken' } );
    }

{{#if signupEnabled}}
    signUp( data ) {
        return this.api.request( {
            url: '/members',
            method: 'post',
            data
        } );
    }

{{/if}}
    login( data ) {
        return this.api.request( {
            url: '/login',
            method: 'post',
            data
        } );
    }

    logout() {
        return this.api.request( {
            url: '/logout',
            method: 'post'
        } );
    }

    updateMemberName( data ) {
        return this.api.request( {
            url: '/members/name',
            method: 'put',
            data
        } );
    }

    updateMemberEmail( data ) {
        return this.api.request( {
            url: '/members/email',
            method: 'put',
            data
        } );
    }

    updateMemberPassword( data ) {
        return this.api.request( {
            url: '/members/password',
            method: 'put',
            data
        } );
    }
}
