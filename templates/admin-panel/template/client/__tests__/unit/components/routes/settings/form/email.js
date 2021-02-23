import { DI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormEmail  from '../../../../../../src/components/routes/settings/form/email.sfc';
import renderer from 'sham-ui-test-helpers';

afterEach( () => {
    DI.bind( 'store', null );
} );

it( 'renders correctly', () => {
    const meta = renderer( RoutesSettingsFormEmail, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    expect.assertions( 2 );

    const updateMock = jest.fn();
    DI.bind( 'store', {
        updateMemberEmail: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormEmail, {
        directives: {
            ...directives
        }
    } );

    const formData = {
        email1: 'admin1@gmail.com',
        email2: 'admin1@gmail.com'
    };
    const { component } = meta;
    component.container.querySelector( '[name="email1"]' ).value = formData.email1;
    component.container.querySelector( '[name="email2"]' ).value = formData.email2;
    component.container.querySelector( '[type="submit"]' ).click();
    component.container.querySelector( '[data-test-modal] [data-test-ok-button]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
