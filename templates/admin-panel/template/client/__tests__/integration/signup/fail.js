import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail sign up', async() => {
    expect.assertions( 5 );

    const DI = setup();

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {}, 401 )
        .use( 'post', '/members', {
            Status: 'Member Not Created',
            Messages: [ 'Name must not be empty.' ]
        }, 400, { 'x-csrf-token': axios.defaultMocksData.csrfToken } );

    history.pushState( {}, '', 'http://client.example.com/signup/' );

    await app.start( DI );

    const formData = {
        name: '',
        email: 'admin@gmail.com',
        password: 'password',
        password2: 'password'
    };
    app.form.fill( 'name', formData.name );
    app.form.fill( 'email', formData.email );
    app.form.fill( 'password', formData.password );
    app.form.fill( 'password2', formData.password2 );
    await app.form.submit();

    app.checkBody();
    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/members' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    expect( window.location.href ).toBe( 'http://client.example.com/signup/' );
} );
