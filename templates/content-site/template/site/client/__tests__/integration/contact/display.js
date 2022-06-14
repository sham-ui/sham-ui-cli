import setup, { app, setupSSR } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();
    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/contact' );

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/contact' );
} );

it( 'ssr', async() => {
    expect.assertions( 1 );
    axios.useDefaultMocks();
    setupSSR();
    await app.ssrAndCheck( 'http://client.example.com/contact' );
} );

it( 'ssr & rehydrate', async() => {
    expect.assertions( 2 );
    axios.useDefaultMocks();

    await app.ssrAndRehydrate( 'http://client.example.com/contact' );

    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/contact' );
} );
