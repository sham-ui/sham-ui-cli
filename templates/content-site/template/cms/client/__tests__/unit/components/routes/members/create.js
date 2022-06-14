import * as directives from 'sham-ui-directives';
import RoutesMembersCreate  from '../../../../../src/components/routes/members/create.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesMembersCreate, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
