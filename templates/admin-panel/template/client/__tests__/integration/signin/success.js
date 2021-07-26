import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success login', async() => {
    expect.assertions( 12 );

    const DI = setup();

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {}, 401 )
        .use( 'post', '/login', {
            Name: 'Admin',
            Email: 'admin@gmail.com'
        }, 200, { 'x-csrf-token': axios.defaultMocksData.csrfToken } );

    history.pushState( {}, '', 'http://client.example.com/login/' );

    await app.start( DI );
    app.checkBody();

    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 0 ][ 0 ] ).toBe( '/csrftoken' );
    expect( axios.mocks.get.mock.calls[ 1 ][ 0 ] ).toBe( '/validsession' );

    axios.use( 'get', '/validsession', {
        Name: 'Admin',
        Email: 'admin@gmail.com'
    } );

    const formData = {
        email: 'admin@gmail.com',
        password: 'password'
    };
    app.form.fill( 'email', formData.email );
    app.form.fill( 'password', formData.password );
    await app.form.submit();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/login' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );

    expect( window.location.href ).toBe( 'http://client.example.com/' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 3 );
    expect( axios.mocks.get.mock.calls[ 1 ][ 0 ] ).toBe( '/validsession' );
    expect( DI.resolve( 'session:storage' ).isAuthenticated ).toBe( true );

    app.checkBody();
} );
