import { DI } from 'sham-ui';
import Router from 'sham-ui-router';
import SignupPage from '../components/routes/signup/page.sfc';
import LoginPage from '../components/routes/login/page.sfc';
import HomePage from '../components/routes/home/page.sfc';
import SettingsPage from '../components/routes/settings/page.sfc';

export default function() {
    const router = new Router( document.location.origin + '/' );

    // Cached home page & login URL
    let homePageURL;
    let loginPageURL;

    router
        .bindPage( '/signup', 'signup', SignupPage, {} )
        .bindPage( '/login', 'login', LoginPage, {} )
        .bindPage( '/settings', 'settings', SettingsPage, {} )
        .bindLazyPage(
            '/members',
            'members/list',
            () => import(
                /* webpackChunkName: "su_members_list" */ '../components/routes/members/page.sfc'
            ),
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
                session.validateSessionPromise.then( isAuthenticated => {
                    if ( [ 'signup', 'login' ].includes( currentRoute.name ) ) {
                        done( !isAuthenticated );
                        if ( isAuthenticated ) {

                            // Authenticated member can't visit signup & login page
                            router.navigate( homePageURL );
                        } else {
                            routerResolve();
                        }
                    } else if (
                        currentRoute.name.startsWith( 'members/' ) &&
                        !session.data.isSuperuser
                    ) {

                        // 403 page
                        done( false );
                        router.navigate( homePageURL );
                    } else {
                        done( isAuthenticated );
                        if ( isAuthenticated ) {
                            routerResolve();
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

function routerResolve() {
    const storage = DI.resolve( 'app:storage' );
    storage.routerResolved = true;
    storage.sync();
}
