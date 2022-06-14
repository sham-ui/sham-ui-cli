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

it( 'menu for small screen size', async() => {
    expect.assertions( 3 );

    axios.useDefaultMocks();
    const DI = setup();

    window.matchMedia.mockImplementation(
        () => ( {
            addListener: jest.fn(),
            matches: false
        } )
    );

    // Open
    await app.start( DI );
    app.click( '.icon-menu' );
    await app.waitRendering();
    app.checkBody();

    // Close
    app.click( '.icon-menu' );
    await app.waitRendering();
    app.checkBody();

    // Open
    app.click( '.icon-menu' );
    await app.waitRendering();

    // Close on wrapper
    app.click( '.show-left' );
    await app.waitRendering();
    app.checkBody();
} );
