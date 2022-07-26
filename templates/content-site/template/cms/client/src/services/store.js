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
            'http://localhost:3003/api/';
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

    articleCategories() {
        return this.api.request( { url: '/categories' } );
    }

    createArticlesCategory( data ) {
        return this.api.request( {
            url: '/categories',
            method: 'post',
            data
        } );
    }

    updateArticlesCategory( id, data ) {
        return this.api.request( {
            url: `/categories/${id}`,
            method: 'put',
            data
        } );
    }

    deleteArticleCategory( id ) {
        return this.api.request( {
            url: `/categories/${id}`,
            method: 'delete'
        } );
    }

    articleTags() {
        return this.api.request( { url: '/tags' } );
    }

    createArticlesTag( data ) {
        return this.api.request( {
            url: '/tags',
            method: 'post',
            data
        } );
    }

    updateArticlesTag( id, data ) {
        return this.api.request( {
            url: `/tags/${id}`,
            method: 'put',
            data
        } );
    }

    deleteArticleTag( id ) {
        return this.api.request( {
            url: `/tags/${id}`,
            method: 'delete'
        } );
    }

    articles( offset, limit ) {
        return this.api.request( {
            url: '/articles',
            params: {
                offset,
                limit
            }
        } );
    }

    articleDetail( id ) {
        return this.api.request( {
            url: `/articles/${id}`
        } );
    }

    createArticle( data ) {
        return this.api.request( {
            url: '/articles',
            method: 'post',
            data
        } );
    }

    updateArticle( id, data ) {
        return this.api.request( {
            url: `/articles/${id}`,
            method: 'put',
            data
        } );
    }

    deleteArticle( id ) {
        return this.api.request( {
            url: `/articles/${id}`,
            method: 'delete'
        } );
    }
}
