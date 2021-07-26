import Router from 'sham-ui-router';
{{#if signupEnabled}}
import SignupPage from '../components/routes/signup/page.sfc';
{{/if}}
import LoginPage from '../components/routes/login/page.sfc';
import HomePage from '../components/routes/home/page.sfc';

export default function( DI ) {
    const router = new Router( DI, document.location.origin + '/' );

    // Cached home page & login URL
    let homePageURL;
    let loginPageURL;

    router
        {{#if signupEnabled}}
        .bindPage( '/signup', 'signup', SignupPage, {} )
        {{/if}}
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
                    if ( [ {{#if signupEnabled}}'signup', {{/if}}'login' ].includes( currentRoute.name ) ) {
                        done( !isAuthenticated );
                        if ( isAuthenticated ) {

                            // Authenticated member can't visit {{#if signupEnabled}}signup &{{/if}}login page
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
                            router.navigate( loginPageURL );
                        }
                    }
                } );
            }
        } );

    homePageURL = router.generate( 'home' );
    loginPageURL = router.generate( 'login' );
    router.resolve();
}

function routerResolve( DI ) {
    const storage = DI.resolve( 'app:storage' );
    storage.routerResolved = true;
    storage.sync();
}
