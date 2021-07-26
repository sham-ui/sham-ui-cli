import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    window.matchMedia = jest.fn().mockImplementation(
        () => ( {
            addListener: jest.fn(),
            matches: true
        } )
    );
} );

afterEach( () => {
    delete window.matchMedia;
} );

it( 'menu', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();

    const DI = setup();
    await app.start( DI );
    app.click( '.icon-menu' );
    await app.waitRendering();
    app.checkBody();

    app.click( '.icon-menu' );
    await app.waitRendering();
    app.checkBody();
} );
