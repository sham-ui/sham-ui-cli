import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success create member', async() => {
    expect.assertions( 5 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [],
            meta: {
                total: 0,
                offset: 0,
                limit: 50
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        name: 'John Smith',
        email: 'john.smith@test.com',
        is_superuser: false,
        password: 'test'
    };

    axios
        .use( 'post', 'admin/members', {} )
        .use( 'get', 'admin/members', {
            members: [
                {
                    ID: 1,
                    Name: formData.name,
                    Email: formData.email,
                    IsSuperuser: formData.is_superuser
                }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } );

    app.form.fill( 'name', formData.name );
    app.form.fill( 'email', formData.email );
    app.form.fill( 'pass', formData.password );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( 'admin/members' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
} );

it( 'success create superuser member', async() => {
    expect.assertions( 5 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [],
            meta: {
                total: 0,
                offset: 0,
                limit: 50
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start();

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        name: 'John Smith',
        email: 'john.smith@test.com',
        is_superuser: true,
        password: 'test'
    };

    axios
        .use( 'post', 'admin/members', {} )
        .use( 'get', 'admin/members', {
            members: [
                {
                    ID: 1,
                    Name: formData.name,
                    Email: formData.email,
                    IsSuperuser: formData.is_superuser
                }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } );

    app.form.fill( 'name', formData.name );
    app.form.fill( 'email', formData.email );
    app.click( '[name="is_superuser"]' );
    app.form.fill( 'pass', formData.password );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( 'admin/members' );
    expect( axios.mocks.post.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
} );
