import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import RoutesSettingsPage  from '../../../../../src/components/routes/settings/page.sfc';
import { storage } from '../../../../../src/storages/session';
import Session from '../../../../../src/services/session';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    const session = storage( DI );
    session.name = 'Test member';
    session.email = 'test@test.com';
    session.sessionValidated = true;
    new Session( DI );

    const meta = renderer( RoutesSettingsPage, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
