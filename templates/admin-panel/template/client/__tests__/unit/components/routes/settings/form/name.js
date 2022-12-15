import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormName  from '../../../../../../src/components/routes/settings/form/name.sfc';
import renderer from 'sham-ui-test-helpers';
import { createDI } from 'sham-ui';

it( 'renders correctly', () => {
    const DI = createDI();
    const meta = renderer( RoutesSettingsFormName, {}, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    expect.assertions( 2 );

    const DI = createDI();
    const updateMock = jest.fn();
    DI.bind( 'store', {
        updateMemberName: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormName, {}, {
        DI,
        directives
    } );

    const formData = {
        newName: 'Johny Smithy'
    };
    const { ctx } = meta;
    ctx.container.querySelector( '[name="name"]' ).value = formData.newName;
    ctx.container.querySelector( '[type="submit"]' ).click();
    ctx.container.querySelector( '[data-test-modal] [data-test-ok-button]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
