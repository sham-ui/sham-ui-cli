import escapeHTML from 'escape-html';
import { createDI, start } from 'sham-ui';
import { setup as setupHydrator, hydrate } from 'sham-ui-ssr/lib/hydrator';
import mainInitializer from './initializers/main';

/**
 * Render application without DOM & any browser environment
 * @param {string} apiURL
 * @param {string} origin Value for router document.location.origin
 * @param {string} href Value for router current location
 * @param {string} cookie Original request cookies
 * @param {Object<string,string>} headers optional headers
 * @return {Promise<{title: string, html: string, data: string, darkThemeEnabled: boolean}>}
 */
export function renderAPP( apiURL, origin, href, cookie, headers ) {
    const DI = createDI();
    DI
        .bind( 'api:url', apiURL )
        .bind( 'api:headers', headers )
        .bind( 'location:origin', origin )
        .bind( 'location:href', href )
        .bind( 'document', { // Container for services
            cookie,
            title: '',
            content: ''
        } )
    ;

    return mainInitializer( DI, setupHydrator( DI ) )
        .then( () => start( DI ) )
        .then( () => hydrate( DI ) )
        .then( storage => {
            const doc = DI.resolve( 'document' );
            const data = {
                ...storage.hydrate(),
                darkThemeEnabled: DI.resolve( 'app:storage' ).darkThemeEnabled,
                title: escapeHTML( doc.title ),
                content: escapeHTML( doc.content )
            };
            DI.resolve( 'router' ).destroy();
            return data;
        } );
}

/**
 * Build result HTML to send go-backend
 * @param {string} data Data for rehydration
 * @param {string} title
 * @param {string} html
 * @param {boolean} darkThemeEnabled
 * @param {string} content
 * @return {string}
 */
export function toHTML( { data, title, html, darkThemeEnabled, content } ) {
    const darkTheme = darkThemeEnabled ?
        ' dark' :
        ''
    ;
    const metaDescription = content ?
        `<meta name="description" content="${content}"/>` :
        ''
    ;
    return (
        '<!doctype html><html><head><meta charset="UTF-8">' +
        '<meta name="viewport" content="width=device-width, initial-scale=1.0">' +
        `<meta name="theme-color" content="${darkThemeEnabled ? '#2b2b2b' : ''}">` +
        metaDescription +
        '<link rel="icon" href="/favicon.ico">' +
        `<title>${title}</title><link rel="stylesheet" href="/bundle.css" /></head>` +
        `<body class="animation-stopped${darkTheme}">${html}<script>window.data=${data};</script>` +
        '<script src="/s.min.js"></script>' +
        '<script>System.import(\'/bundle.js\');</script></body></html>'
    );
}

process.on( 'message', function( msg ) {
    renderAPP( msg.api, msg.origin, msg.url, msg.cookies, msg.headers )
        .then(
            data => ( {
                id: msg.id,
                html: toHTML( data )
            } ),

            // Catch & wrap any error for process in go-backend
            err => ( {
                id: msg.id,
                error: `${err.toString()}\n${err.stack}`
            } )
        ).then(
            data => process.send( data )
        );
} );
