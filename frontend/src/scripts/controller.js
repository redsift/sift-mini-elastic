/**
 * Mini Elastic Sift. Frontend controller entry point.
 */
import { SiftController, registerSiftController } from '@redsift/sift-sdk-web';

export default class MyController extends SiftController {
  constructor() {
    // You have to call the super() method to initialize the base class.
    super();
    this._suHandler = this.onStorageUpdate.bind(this);
  }

  // for more info: http://docs.redsift.com/docs/client-code-siftcontroller
  loadView(state) {
    console.log('mini-elastic: loadView', state);
    // Register for storage update events on the "stats" bucket so we can update the UI
    this.storage.subscribe(['stats'], this._suHandler);
    const apiUrl = state.params.rpcApiConfig.baseUrl
    switch (state.type) {
      case 'summary':
        return {
          html: 'summary.html',
          data: this.getStats().then(stats => ({stats, apiUrl}))
        };
      default:
        console.error('mini-elastic: unknown Sift type: ', state.type);
    }
  }

  // Event: storage update
  onStorageUpdate(value) {
    console.log('mini-elastic: onStorageUpdate: ', value);
    return this.getStats().then(xe => {
      // Publish events from 'stats' to view
      this.publish('stats', xe);
    });
  }

  getStats() {
    return this.storage.getAll({
      bucket: 'stats'
    }).then(v => {
      console.log('mini-elastic: getStats returned:', v);
      return v.length > 0 ? JSON.parse(v[0].value) : {}
    });
  }

}

// Do not remove. The Sift is responsible for registering its views and controllers
registerSiftController(new MyController());
