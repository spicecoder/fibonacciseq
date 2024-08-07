function initializeSequence(pnr) {
    pnr.fibonacci = [0, 1];
    //possible Object : Starter ,intention :setup first two members in the seq

  }
  
  function addTwoNumbers(pnr) {
    const seq = pnr.fibonacci;
    const lastIndex = seq.length - 1;
    return seq[lastIndex] + seq[lastIndex - 1];
     //possible Object : Joiner  ,intention :add  two members 

  }
  
  function repeatAddition(pnr) {
    let nextFib = addTwoNumbers(pnr);
    while (nextFib <= pnr.max) {
      pnr.fibonacci.push(nextFib);
      nextFib = addTwoNumbers(pnr);
    }
    // possible Object : NextFib ,calculateNext  :calculate next in the seq

  }
  
  function filterSequence(pnr) {
    pnr.fibonacciInRange = pnr.fibonacci.filter(num => num >= pnr.min && num <= pnr.max);
  // possible Object : GoNotGo ,Checkcondition :calculate if request met

}
  
  function generateFibonacciInRange(min, max) {
    if (min > max) {
      return "Invalid range";
    }
  
    const pnr = { min, max, fibonacci: [], fibonacciInRange: [],intention_satisfied:'N' };
  
    // Step 1: Initialize the Sequence
    initializeSequence(pnr);
  
    // Step 3: Repeat the Addition
    repeatAddition(pnr);
  
    // Step 4: Filter the Sequence
    filterSequence(pnr);
  
    return pnr.fibonacciInRange;
  }
  
  // Example usage
  const min = 5;
  const max = 55;
  const fibonacciInRange = generateFibonacciInRange(min, max);
  console.log(fibonacciInRange); // Output: [5, 8, 13, 21, 34, 55]
  /*
dc : human, emits intention "setup first two members in the seq" with PnR :{min,max} ->object: Starter:reflects to 
initializeSequence(dc),which in turn emits intention :"find next fb" with PNR: {currentsrq}->object->fbSequence ->reflects-"calculate next memnber"-calcNextseq-emits :"set fbseq"-to
object fbsequence -reflects itention checkResultCondition to EvaluateFibomacciReach 
PnR ={"setup done":"Y/N", currentseq:"object:fbsequence",}
  */
  