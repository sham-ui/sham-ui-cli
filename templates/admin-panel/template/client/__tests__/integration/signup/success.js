import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success sign up', async() => {
    expect.assertions( 15 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {}, 401 )
        .use( 'post', '/members', {
            Status: 'Member Created',
            Messages: null
        }, 200, { 'x-csrf-token': axios.defaultMocksData.csrfToken } );

    history.pushState( {}, '', 'http://client.example.com/signup/' );

    await app.start( DI );
    app.checkBody();

    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 0 ][ 0 ] ).toBe( '/csrftoken' );
    expect( axios.mocks.get.mock.calls[ 1 ][ 0 ] ).toBe( '/validsession' );

    const formData = {
        name: 'admin',
        email: 'admin@gmail.com',
        password: 'password',
        password2: 'password'
    };
    app.form.fill( 'name', formData.name );
    app.form.fill( 'email', formData.email );
    app.form.fill( 'password', formData.password );
    app.form.fill( 'password2', formData.password2 );
    await app.form.submit();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/members' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    expect( axios.mocks.get.mock.calls[ 2 ][ 0 ] ).toBe( '/validsession' );

    expect( window.location.href ).toBe( 'http://client.example.com/login' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 3 );

    axios
        .use( 'get', '/validsession', {
            Name: 'Admin',
            Email: 'admin@gmail.com'
        } )
        .use( 'post', '/login', {
            Name: 'Admin',
            Email: 'admin@gmail.com'
        }, 200, { 'x-csrf-token': axios.defaultMocksData.csrfToken } );

    await app.waitRendering();
    app.form.fill( 'email', formData.email );
    app.form.fill( 'password', formData.password );
    await app.form.submit();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 4 );
    expect( axios.mocks.get.mock.calls[ 3 ][ 0 ] ).toBe( '/validsession' );
    expect( DI.resolve( 'session:storage' ).isAuthenticated ).toBe( true );

    app.checkBody();
} );
