import Router from 'sham-ui-router';
import LoginPage from '../components/routes/login/page.sfc';
import HomePage from '../components/routes/home/page.sfc';

export default function( DI ) {
    const router = new Router( DI, document.location.origin + '/' );

    // Cached home page URL
    let homePageURL;

    router
        .bindPage( '/login', 'login', LoginPage, {} )
        .bindLazyPage(
            '/settings',
            'settings',
            () => import( '../components/routes/settings/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/members',
            'members/list',
            () => import( '../components/routes/members/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/server',
            'server-info',
            () => import( '../components/routes/server-info/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/articles/categories',
            'articles/categories',
            () => import( '../components/routes/articles/categories/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/articles/tags',
            'articles/tags',
            () => import( '../components/routes/articles/tags/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/articles/new',
            'articles/new',
            () => import( '../components/routes/articles/new/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/articles/:id/edit',
            'articles/edit',
            () => import( '../components/routes/articles/edit/page.sfc' ),
            {}
        )
        .bindLazyPage(
            '/articles',
            'articles/list',
            () => import( '../components/routes/articles/list/page.sfc' ),
            {}
        )
        .bindPage( '', 'home', HomePage, {} )
        .hooks( {
            before( done ) {
                const currentRoute = router.storage;
                if ( 'home' === currentRoute.name && currentRoute.url !== homePageURL ) {

                    // 404 page
                    done( false );
                    router.navigate( homePageURL );
                    return;
                }
                const session = DI.resolve( 'session' );
                session.validateSession().then( isAuthenticated => {
                    if ( [ 'login' ].includes( currentRoute.name ) ) {
                        done( !isAuthenticated );
                        if ( isAuthenticated ) {

                            // Authenticated member can't visit login page
                            router.navigate( homePageURL );
                        } else {
                            routerResolve( DI );
                        }
                    } else if (
                        (
                            currentRoute.name.startsWith( 'members/' ) ||
                            'server-info' === currentRoute.name
                        ) &&
                        !session.data.isSuperuser
                    ) {

                        // 403 page
                        done( false );
                        router.navigate( homePageURL );
                    } else {
                        done( isAuthenticated );
                        if ( isAuthenticated ) {
                            routerResolve( DI );
                        } else {

                            // If non authenticated then redirects to login
                            router.navigateToRoute( 'login' );
                        }
                    }
                } );
            }
        } );

    homePageURL = router.generate( 'home' );
    router.resolve();
}

function routerResolve( DI ) {
    const storage = DI.resolve( 'app:storage' );
    storage.routerResolved = true;
    storage.sync();
}
