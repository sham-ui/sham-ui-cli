import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );
const DI = setup();

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success reset member password', async() => {
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
    await app.start( DI );

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.checkMainPanel();

    const formData = {
        pass1: 'password',
        pass2: 'password'
    };

    app.form.fill( 'pass1', formData.pass1 );
    app.form.fill( 'pass2', formData.pass2 );
    app.click( '.form-layout:last-child [type="submit"]' );
    await app.waitRendering();

    axios.use( 'put', 'admin/members/2/password', {} );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( 'admin/members/2/password' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
    app.checkMainPanel();
} );
