import AssetsPanel  from '../../../../../src/components/routes/server-info/assets-panel.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( AssetsPanel, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
