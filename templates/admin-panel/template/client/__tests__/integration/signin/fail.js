import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail login', async() => {
    expect.assertions( 3 );

    const DI = setup();

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {}, 401 )
        .use( 'post', '/login', {
            Status: 'Failed to authenticate',
            Messages: [ 'Incorrect username or password' ]
        }, 401 );

    history.pushState( {}, '', 'http://client.example.com/login/' );

    await app.start( DI );

    const formData = {
        email: 'admin@gmail.com',
        password: 'password'
    };
    app.form.fill( 'email', formData.email );
    app.form.fill( 'password', formData.password );
    await app.form.submit();

    expect( window.location.href ).toBe( 'http://client.example.com/login/' );
    expect( DI.resolve( 'session:storage' ).isAuthenticated ).toBe( false );

    app.checkBody();
} );
