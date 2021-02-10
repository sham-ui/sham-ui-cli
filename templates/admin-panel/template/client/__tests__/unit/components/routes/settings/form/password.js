import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormPassword  from '../../../../../../src/components/routes/settings/form/password.sfc';
import renderer from 'sham-ui-test-helpers';

afterEach( () => {
    DI.bind( 'store', null );
} );

it( 'renders correctly', () => {
    const meta = renderer( RoutesSettingsFormPassword, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    expect.assertions( 2 );

    const updateMock = jest.fn();
    DI.bind( 'store', {
        updateMemberPassword: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormPassword, {
        directives
    } );

    const formData = {
        pass1: 'admin1@gmail.com',
        pass2: 'admin1@gmail.com'
    };
    const { component } = meta;
    component.container.querySelector( '[name="pass1"]' ).value = formData.pass1;
    component.container.querySelector( '[name="pass2"]' ).value = formData.pass2;
    component.container.querySelector( '[type="submit"]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
