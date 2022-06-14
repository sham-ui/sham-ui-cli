import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app, setupSSR } from '../helpers';
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
        .use( 'get', '/articles/pervya-testovaya-statya', { Messages: [ 'Fail load' ] }, 500 )
    ;

    history.pushState( {}, '', 'http://client.example.com/pervya-testovaya-statya' );

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/pervya-testovaya-statya' );
} );


it( 'ssr', async() => {
    expect.assertions( 1 );
    axios
        .useDefaultMocks()
        .use( 'get', '/articles/pervya-testovaya-statya', { Messages: [ 'Fail load' ] }, 500 )
    ;
    setupSSR();
    await app.ssrAndCheck( 'http://client.example.com/pervya-testovaya-statya' );
} );


it( 'ssr & rehydrate', async() => {
    expect.assertions( 2 );

    axios
        .useDefaultMocks()
        .use( 'get', '/articles/pervya-testovaya-statya', { Messages: [ 'Fail load' ] }, 500 )
    ;

    await app.ssrAndRehydrate( 'http://client.example.com/pervya-testovaya-statya' );

    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/pervya-testovaya-statya' );
} );
