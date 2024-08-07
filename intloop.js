// Initialize the PnR set
const pnr = {
    Fbsequence: {
      started: 'N',
      sequence: [],
      genComplete: 'N'
    },
    min: 5,
    max: 55,
    fibonacciInRange: []
  };
  
  // Define Design Chunks
  const designChunks = {
    startSequence: function(pnr) {
      if (pnr.Fbsequence.started === 'N') {
        const intention = "Setup first two members in the fbsequence";
        pnr.Fbsequence.sequence = [0, 1];
        pnr.Fbsequence.started = 'Y';
        processIntention(intention, pnr);
      }
    },
    findNextFibonacci: function(pnr) {
      if (pnr.Fbsequence.started === 'Y' && pnr.Fbsequence.genComplete === 'N') {
        const nextFib = addTwoNumbers(pnr.Fbsequence.sequence);
        if (nextFib > pnr.max) {
          pnr.Fbsequence.genComplete = 'Y';
        } else {
          const intention = "Set Fibonacci sequence";
          pnr.Fbsequence.sequence.push(nextFib);
          processIntention(intention, pnr);
        }
      }
    },
    filterSequence: function(pnr) {
      if (pnr.Fbsequence.genComplete === 'Y') {
        pnr.fibonacciInRange = pnr.Fbsequence.sequence.filter(num => num >= pnr.min && num <= pnr.max);
      }
    }
  };
  
  // Process Intentions
  function processIntention(intention, pnr) {
    switch (intention) {
      case "Setup first two members in the fbsequence":
        // Already handled in startSequence
        break;
      case "Find next Fibonacci sequence":
        // Already handled in findNextFibonacci
        break;
      case "Set Fibonacci sequence":
        // No additional action needed
        break;
    }
  }
  
  // Add Two Numbers Function
  function addTwoNumbers(sequence) {
    const lastIndex = sequence.length - 1;
    return sequence[lastIndex] + sequence[lastIndex - 1];
  }
  
  // Intention Ring Loop
  function intentionRing(pnr) {
    let executed;
    do {
      executed = false;
      if (pnr.Fbsequence.started === 'N') {
        designChunks.startSequence(pnr);
        executed = true;
      }
      if (pnr.Fbsequence.started === 'Y' && pnr.Fbsequence.genComplete === 'N') {
        designChunks.findNextFibonacci(pnr);
        executed = true;
      }
      if (pnr.Fbsequence.genComplete === 'Y' && !pnr.fibonacciInRange.length) {
        designChunks.filterSequence(pnr);
        executed = true;
      }
    } while (executed);
  }
  
  // Execute the Intention Ring Loop
  intentionRing(pnr);
  console.log(pnr.fibonacciInRange); // Output: [5, 8, 13, 21, 34, 55]
  