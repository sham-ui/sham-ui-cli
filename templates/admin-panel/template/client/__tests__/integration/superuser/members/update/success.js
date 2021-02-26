import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'success update member data', async() => {
    expect.assertions( 5 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 2, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start();

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        name: 'New name',
        email: 'new.john.smith.1@test.com',
        is_superuser: true
    };

    app.form.fill( 'name', formData.name );
    app.form.fill( 'email', formData.email );
    app.click( '[name="is_superuser"]' );
    app.form.submit();

    axios
        .use( 'put', 'admin/members/2', {} )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 2, Name: formData.name, Email: formData.email, IsSuperuser: formData.email }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( 'admin/members/2' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkMainPanel();
} );
