import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'logout from settings page', async() => {
    expect.assertions( 7 );

    axios.useDefaultMocks();
    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start( DI );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );

    app.click( '.panel.settings p:nth-of-type(1) .icon-pencil' );

    const formData = {
        name: ''
    };
    axios
        .use( 'put', '/members/name', {
            'Status': 'Bad Name',
            'Messages': [ 'Name must have more than 0 characters.' ]
        }, 401 )
        .use( 'post', '/logout' )
        .use( 'get', '/validsession', {}, 401 );

    app.form.fill( 'name', formData.name );
    await app.form.submit();
    app.click( '[data-test-modal] [data-test-ok-button]' );

    // Wait request handling
    await app.waitRendering();

    // Wait requestAnimationFrame & logout
    await app.waitRendering();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/logout' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 3 );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/login' );
} );
