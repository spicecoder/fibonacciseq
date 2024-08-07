// Class for Intention
class Intention {
    constructor(name, subset) {
      this.name = name;
      this.subset = subset;
    }
  }
  
  // Class for Object (Fbsequence)
  class Fbsequence {
    constructor() {
      this.started = 'N';
      this.sequence = [];
      this.genComplete = 'N';
    }
  
    receiveIntention(intention, pnr) {
      switch (intention.name) {
        case "Setup first two members in the fbsequence":
          this.sequence = intention.subset.sequence;
          this.started = intention.subset.started;
          this.reflectIntention("Find next Fibonacci sequence", pnr);
          break;
        case "Set Fibonacci sequence":
          this.sequence.push(intention.subset.nextFib);
          if (intention.subset.nextFib >= pnr.min && intention.subset.nextFib <= pnr.max) {
            pnr.fibonacciInRange.push(intention.subset.nextFib);
          }
          this.reflectIntention("Find next Fibonacci sequence", pnr);
          break;
      }
    }
  
    reflectIntention(intentionName, pnr) {
      const designChunk = pnr.designChunks.find(chunk => chunk.name === intentionName);
      if (designChunk) {
        designChunk.absorbIntention(pnr);
      }
    }
  }
  
  // Class for DesignChunk
  class DesignChunk {
    constructor(name, action) {
      this.name = name;
      this.action = action;
    }
  
    emitIntention(intention, pnr) {
      pnr.fbsequence.receiveIntention(intention, pnr);
    }
  
    absorbIntention(pnr) {
      this.action(pnr);
    }
  }
  
  // Class for Intention Ring
  class IntentionRing {
    constructor(pnr) {
      this.pnr = pnr;
    }
  
    execute() {
      let executed;
      do {
        executed = false;
        for (const chunk of this.pnr.designChunks) {
          if (chunk.name === "Setup first two members in the fbsequence" && this.pnr.fbsequence.started === 'N') {
            chunk.absorbIntention(this.pnr);
            executed = true;
          } else if (chunk.name === "Find next Fibonacci sequence" && this.pnr.fbsequence.started === 'Y' && this.pnr.fbsequence.genComplete === 'N') {
            chunk.absorbIntention(this.pnr);
            executed = true;
          }
        }
      } while (executed);
    }
  }
  
  // Add Two Numbers Function
  function addTwoNumbers(sequence) {
    const lastIndex = sequence.length - 1;
    return sequence[lastIndex] + sequence[lastIndex - 1];
  }
  
  // PnR Set Initialization
  const pnr = {
    fbsequence: new Fbsequence(),
    min: 5,
    max: 55,
    fibonacciInRange: [],
    designChunks: []
  };
  
  // Define Design Chunks
  const startSequenceChunk = new DesignChunk("Setup first two members in the fbsequence", (pnr) => {
    if (pnr.fbsequence.started === 'N') {
      const intention = new Intention("Setup first two members in the fbsequence", {
        sequence: [0, 1],
        started: 'Y'
      });
      startSequenceChunk.emitIntention(intention, pnr);
    }
  });
  
  const findNextFibonacciChunk = new DesignChunk("Find next Fibonacci sequence", (pnr) => {
    if (pnr.fbsequence.started === 'Y' && pnr.fbsequence.genComplete === 'N') {
      const nextFib = addTwoNumbers(pnr.fbsequence.sequence);
      if (nextFib > pnr.max) {
        pnr.fbsequence.genComplete = 'Y';
      } else {
        const intention = new Intention("Set Fibonacci sequence", {
          nextFib: nextFib
        });
        findNextFibonacciChunk.emitIntention(intention, pnr);
      }
    }
  });
  
  // Register Design Chunks
  pnr.designChunks.push(startSequenceChunk, findNextFibonacciChunk);
  
  // Create and Execute the Intention Ring Loop
  const intentionRing = new IntentionRing(pnr);
  intentionRing.execute();
  console.log(pnr.fibonacciInRange); // Output: [5, 8, 13, 21, 34, 55]
  