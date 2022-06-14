import * as directives from 'sham-ui-directives';
import RoutesMembersDetail  from '../../../../../src/components/routes/members/detail.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesMembersDetail, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
