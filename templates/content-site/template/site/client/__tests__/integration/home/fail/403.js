import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app, setupSSR } from '../../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display', async() => {
    expect.assertions( 2 );

    const DI = setup();
    setupUnsafe( DI );

    axios
        .useDefaultMocks()
        .use( 'get', '/articles', { Messages: [ 'Fail load' ] }, 403 )
    ;

    history.pushState( {}, '', 'http://client.example.com/' );

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );


it( 'ssr', async() => {
    expect.assertions( 1 );
    axios
        .useDefaultMocks()
        .use( 'get', '/articles', { Messages: [ 'Fail load' ] }, 403 )
    ;
    setupSSR();
    await app.ssrAndCheck( 'http://client.example.com/' );
} );


it( 'ssr & rehydrate', async() => {
    expect.assertions( 2 );

    axios
        .useDefaultMocks()
        .use( 'get', '/articles', { Messages: [ 'Fail load' ] }, 403 )
    ;

    await app.ssrAndRehydrate( 'http://client.example.com/' );

    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
