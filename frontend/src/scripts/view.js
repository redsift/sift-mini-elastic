/**
 * Mini Elastic Sift. Frontend view entry point.
 */
import { SiftView, registerSiftView } from '@redsift/sift-sdk-web';

const nanoToSecs = v => Math.round(v/Math.pow(10, 7)) / 100

export default class MyView extends SiftView {
  constructor() {
    // You have to call the super() method to initialize the base class.
    super();
    this.onStorageUpdate = this.onStorageUpdate.bind(this);
  }

  // for more info: http://docs.redsift.com/docs/client-code-siftview
  presentView(value) {
    console.log('mini-elastic: presentView: ', value);
    document.querySelector("#api_url").innerHTML = value.data.apiUrl;
    this.controller.subscribe('storageupdated', this.onStorageUpdate);
    this.onStorageUpdate(value.data)
  };

  willPresentView(value) { };


  onStorageUpdate({stats}) {
    if (!stats || !stats.index) {
      return;
    }
    ["analysis_time", "index_time"].forEach(x => stats.index[x] = nanoToSecs(stats.index[x]) +'s' );
    document.querySelector("#index_stats").innerHTML = JSON.stringify(stats, null, ' ')
  }

}

registerSiftView(new MyView(window));
