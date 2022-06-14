import Router from 'sham-ui-router';
import HomePage from '../components/routes/home/page.sfc';
import ArticlePage from '../components/routes/article/page.sfc';

export default function( DI ) {
    const router = new Router( DI, DI.resolve( 'location:origin' ) );

    router
        .bindPage( '/page/:page/', 'home.page', HomePage, {} )
        .bindLazyPage(
            '/category/:category/:page/',
            'category.page',
            () => import( '../components/routes/category/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/tag/:tag/:page/',
            'tag.page',
            () => import( '../components/routes/tag/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/search/:page/',
            'search.page',
            () => import( '../components/routes/search/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/contact/',
            'contact',
            () => import( '../components/routes/contact/page.sfc' ),
            {}
        )
        .bindPage( '/:slug', 'article', ArticlePage, {} )
        .bindPage( '', 'home', HomePage, {} )
    ;

    return new Promise( resolve => {
        router.hooks( {
            after() {
                resolve();
            }
        } );
        router.resolve( DI.resolve( 'location:href' ) );
    } ).then( () => {
        const routerStorage = DI.resolve( 'router:storage' );
        return routerStorage.pageLoaded ?
            null :
            new Promise( resolve => {
                const changed = () => {
                    routerStorage.removeWatcher( 'pageLoaded', changed );
                    resolve();
                };
                routerStorage.addWatcher( 'pageLoaded', changed );
            } );
    } ).then( () => {
        const appStorage = DI.resolve( 'app:storage' );
        appStorage.routerResolved = true;
        appStorage.sync();
    } );
}
